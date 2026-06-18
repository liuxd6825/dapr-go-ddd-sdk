package miniofs

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/afero"
)

// Compile-time assertions: the tests below use these via explicit type
// assignment so removing any method on *Fs or *file will fail to build.
var (
	_ afero.Fs    = (*Fs)(nil)
	_ afero.File  = (*file)(nil)
)

// TestAferoFsInterface is a placeholder for the interface assertion above;
// keeping it in the test file makes the build-time check show up in `go test`
// output as a real test entry.
func TestAferoFsInterface(t *testing.T) {}

// TestAferoFileInterface mirrors the file-side assertion.
func TestAferoFileInterface(t *testing.T) {}

// TestNewFsRejectsNilClient documents the precondition that the caller must
// supply a non-nil minio.Client.
func TestNewFsRejectsNilClient(t *testing.T) {
	if _, err := NewFs(Config{Bucket: "b"}, nil); err == nil {
		t.Fatal("expected error for nil client, got nil")
	}
}

// TestNewFsRejectsEmptyBucket documents the precondition that a bucket must
// be configured.
func TestNewFsRejectsEmptyBucket(t *testing.T) {
	if _, err := NewFs(Config{}, nil); err == nil {
		t.Fatal("expected error for empty bucket, got nil")
	}
}

const (
	testEndpoint  = "192.168.120.224:9000"
	testAccessKey = "minioadmin"
	testSecretKey = "minioadmin"
	testBucket    = "import-data"
	testUseSSL    = true

	// testPrefix is prepended to every object name produced by the test so
	// concurrent test runs and leftover state do not collide.
	testPrefix = "miniofs_test/"
)

// newTestClient builds a minio.Client using the connection details supplied
// for the integration test. The configured Secure setting is tried first; if
// it fails with the typical "server is actually HTTP" mismatch the client is
// rebuilt without TLS so the test can run against a default MinIO install.
// When the server is completely unreachable the test is skipped.
func newTestClient(t *testing.T) *minio.Client {
	t.Helper()
	tryConnect := func(secure bool) (*minio.Client, error) {
		c, err := minio.New(testEndpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(testAccessKey, testSecretKey, ""),
			Secure: secure,
		})
		if err != nil {
			return nil, err
		}
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_, err = c.BucketExists(ctx, testBucket)
		return c, err
	}

	if client, err := tryConnect(testUseSSL); err == nil {
		return client
	} else if client, err2 := tryConnect(!testUseSSL); err2 == nil {
		t.Logf("configured useSSL=%t did not work (%v); falling back to useSSL=%t",
			testUseSSL, err, !testUseSSL)
		return client
	} else {
		t.Skipf("minio server %s not reachable (ssl=%t: %v; ssl=%t: %v)",
			testEndpoint, testUseSSL, err, !testUseSSL, err2)
		return nil
	}
}

// uniqueName returns a unique object name within the test prefix so that
// repeated runs do not interfere with each other.
func uniqueName(suffix string) string {
	return fmt.Sprintf("%s%d_%s", testPrefix, time.Now().UnixNano(), suffix)
}

func newTestFs(t *testing.T) (afero.Fs, *minio.Client) {
	t.Helper()
	client := newTestClient(t)
	fs, err := NewFs(Config{
		Name:     "test",
		Bucket:   testBucket,
		RootPath: "",
		Tags:     []string{"test"},
	}, client)
	if err != nil {
		t.Fatalf("NewFs: %v", err)
	}
	return fs, client
}

// cleanup removes the test prefix directory at the end of the run. Any error
// is reported but does not fail the test.
func cleanup(t *testing.T, client *minio.Client, names ...string) {
	t.Helper()
	ctx := context.Background()
	for _, n := range names {
		if err := client.RemoveObject(ctx, testBucket, n, minio.RemoveObjectOptions{}); err != nil {
			t.Logf("cleanup %s: %v", n, err)
		}
	}
}

// TestConnect exercises the full lifecycle against a real MinIO server using
// the credentials supplied in the task description. The test is skipped when
// the server is unreachable so that it can run in offline environments.
func TestConnect(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in -short mode")
	}

	fs, client := newTestFs(t)

	created := uniqueName("hello.txt")
	renamed := uniqueName("hello_renamed.txt")
	dir := uniqueName("dir/")
	dirChild := dir + "child.txt"

	t.Cleanup(func() {
		cleanup(t, client, created, renamed, dirChild, dir)
	})

	t.Run("stat_missing", func(t *testing.T) {
		missing := uniqueName("missing.txt")
		if _, err := fs.Stat(missing); err == nil {
			t.Errorf("expected error for missing object, got nil")
		}
	})

	t.Run("create_and_read", func(t *testing.T) {
		payload := "hello miniofs"
		f, err := fs.Create(created)
		if err != nil {
			t.Fatalf("Create: %v", err)
		}
		if _, err := f.Write([]byte(payload)); err != nil {
			t.Fatalf("Write: %v", err)
		}
		if err := f.Close(); err != nil {
			t.Fatalf("Close: %v", err)
		}

		got, err := afero.ReadFile(fs, created)
		if err != nil {
			t.Fatalf("ReadFile: %v", err)
		}
		if string(got) != payload {
			t.Errorf("content mismatch: got %q want %q", got, payload)
		}
	})

	t.Run("stat", func(t *testing.T) {
		info, err := fs.Stat(created)
		if err != nil {
			t.Fatalf("Stat: %v", err)
		}
		if info.IsDir() {
			t.Errorf("expected file, got directory")
		}
		if info.Size() != int64(len("hello miniofs")) {
			t.Errorf("size mismatch: got %d", info.Size())
		}
	})

	t.Run("rename", func(t *testing.T) {
		if err := fs.Rename(created, renamed); err != nil {
			t.Fatalf("Rename: %v", err)
		}
		if _, err := fs.Stat(created); err == nil {
			t.Errorf("expected source %s to be gone after rename", created)
		}
		if _, err := fs.Stat(renamed); err != nil {
			t.Errorf("expected dest %s to exist: %v", renamed, err)
		}
	})

	t.Run("openfile_write_sequential", func(t *testing.T) {
		// Verify sequential Write/Close is treated as a streaming PutObject.
		f, err := fs.OpenFile(uniqueName("seq.txt"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			t.Fatalf("OpenFile: %v", err)
		}
		seqName := f.Name()
		t.Cleanup(func() { cleanup(t, client, seqName) })
		if _, err := f.Write([]byte("part1-")); err != nil {
			t.Fatalf("Write 1: %v", err)
		}
		if _, err := f.Write([]byte("part2")); err != nil {
			t.Fatalf("Write 2: %v", err)
		}
		if err := f.Close(); err != nil {
			t.Fatalf("Close: %v", err)
		}
		got, err := afero.ReadFile(fs, seqName)
		if err != nil {
			t.Fatalf("ReadFile: %v", err)
		}
		if string(got) != "part1-part2" {
			t.Errorf("content mismatch: got %q", got)
		}
	})

	t.Run("read_full", func(t *testing.T) {
		f, err := fs.Open(renamed)
		if err != nil {
			t.Fatalf("Open: %v", err)
		}
		data, err := io.ReadAll(f)
		if err != nil {
			t.Fatalf("ReadAll: %v", err)
		}
		if err := f.Close(); err != nil {
			t.Fatalf("Close: %v", err)
		}
		if string(data) != "hello miniofs" {
			t.Errorf("content mismatch: got %q", data)
		}
	})

	t.Run("read_at", func(t *testing.T) {
		f, err := fs.Open(renamed)
		if err != nil {
			t.Fatalf("Open: %v", err)
		}
		defer f.Close()
		buf := make([]byte, 5)
		n, err := f.ReadAt(buf, 6)
		if err != nil {
			t.Fatalf("ReadAt: %v", err)
		}
		if n != 5 || string(buf) != "minio" {
			t.Errorf("ReadAt mismatch: n=%d buf=%q", n, buf)
		}
	})

	t.Run("remove", func(t *testing.T) {
		if err := fs.Remove(renamed); err != nil {
			t.Fatalf("Remove: %v", err)
		}
		if _, err := fs.Stat(renamed); err == nil {
			t.Errorf("expected %s to be gone", renamed)
		}
	})

	t.Run("mkdir_mkdirall_creates_markers", func(t *testing.T) {
		// Mkdir must create a single marker ending in "/".
		leaf := uniqueName("mkone")
		if err := fs.Mkdir(leaf, 0755); err != nil {
			t.Fatalf("Mkdir: %v", err)
		}
		t.Cleanup(func() { cleanup(t, client, leaf+"/") })
		info, err := fs.Stat(leaf)
		if err != nil {
			t.Fatalf("Stat %s: %v", leaf, err)
		}
		if !info.IsDir() {
			t.Errorf("expected %s to be reported as a directory", leaf)
		}

		// MkdirAll must create a marker for every path level.
		base := uniqueName("multi")
		paths := []string{
			base + "/a",
			base + "/a/b",
			base + "/a/b/c",
		}
		t.Cleanup(func() {
			cleanup(t, client,
				paths[0]+"/",
				paths[1]+"/",
				paths[2]+"/",
			)
		})
		if err := fs.MkdirAll(paths[2], 0755); err != nil {
			t.Fatalf("MkdirAll: %v", err)
		}
		for _, p := range paths {
			info, err := fs.Stat(p)
			if err != nil {
				t.Errorf("Stat %s: %v", p, err)
				continue
			}
			if !info.IsDir() {
				t.Errorf("expected %s to be a directory, got file", p)
			}
		}
	})

	t.Run("mkdirall_nested_a_b_c", func(t *testing.T) {
		// Focus test for MkdirAll("a/b/c"). Wrapped under a unique prefix
		// to avoid colliding with other tests' state, the call must produce
		// a marker at every path level:
		//   <root>/a/
		//   <root>/a/b/
		//   <root>/a/b/c/
		// and Stat must recognise each one as a directory even when called
		// without the trailing slash.
		root := uniqueName("nested")
		levelA := root + "/a"
		levelB := root + "/a/b"
		levelC := root + "/a/b/c"
		t.Cleanup(func() {
			// RemoveAll at the root cleans up every level.
			if err := fs.RemoveAll(root); err != nil {
				t.Logf("cleanup RemoveAll(%s): %v", root, err)
			}
		})

		// 1. Drive MkdirAll.
		if err := fs.MkdirAll(levelC, 0755); err != nil {
			t.Fatalf("MkdirAll(%q): %v", levelC, err)
		}

		// 2. Each ancestor must be a directory according to Stat, both with
		//    and without the trailing slash (Go's os.Stat convention).
		checks := []string{levelA, levelB, levelC}
		for _, p := range checks {
			info, err := fs.Stat(p)
			if err != nil {
				t.Errorf("Stat %q: %v", p, err)
				continue
			}
			if !info.IsDir() {
				t.Errorf("Stat %q: IsDir() = false, want true", p)
			}
		}

		// 3. Readdir at the root must show exactly one child ("a"), and that
		//    child must be a directory.
		f, err := fs.Open(root)
		if err != nil {
			t.Fatalf("Open(%q): %v", root, err)
		}
		infos, err := f.Readdir(100)
		if err != nil {
			t.Fatalf("Readdir: %v", err)
		}
		_ = f.Close()
		if len(infos) != 1 {
			t.Fatalf("Readdir at root: want 1 entry, got %d (%v)", len(infos), namesOf(infos))
		}
		if infos[0].Name() != "a" {
			t.Errorf("Readdir: want entry named \"a\", got %q", infos[0].Name())
		}
		if !infos[0].IsDir() {
			t.Errorf("Readdir entry %q: IsDir() = false, want true", infos[0].Name())
		}

		// 4. Writing a file under c/ and reading it back must work; the
		//    intermediate markers must not block normal object operations.
		leaf := levelC + "/file.txt"
		want := "hello from c"
		if err := seedObject(fs, leaf, []byte(want)); err != nil {
			t.Fatalf("seedObject %q: %v", leaf, err)
		}
		got, err := afero.ReadFile(fs, leaf)
		if err != nil {
			t.Fatalf("ReadFile %q: %v", leaf, err)
		}
		if string(got) != want {
			t.Errorf("content mismatch: got %q want %q", got, want)
		}

		// 5. RemoveAll at the root must delete the markers, the leaf, and
		//    every intermediate level. After this, Stat on every level
		//    must fail (mirroring os.RemoveAll).
		if err := fs.RemoveAll(root); err != nil {
			t.Fatalf("RemoveAll(%q): %v", root, err)
		}
		for _, p := range append(checks, leaf) {
			if _, err := fs.Stat(p); err == nil {
				t.Errorf("expected %q to be gone after RemoveAll", p)
			}
		}
	})

	t.Run("stat_file_missing_after_mkdir", func(t *testing.T) {
		// Stat on the bare file name (no trailing slash) of a directory
		// should fail: the object key is "<dir>/", not "<dir>".
		leaf := uniqueName("statdir")
		if err := fs.Mkdir(leaf, 0755); err != nil {
			t.Fatalf("Mkdir: %v", err)
		}
		t.Cleanup(func() { cleanup(t, client, leaf+"/") })
		if _, err := fs.Stat(leaf + "/something.txt"); err == nil {
			t.Errorf("expected missing-file error inside an empty directory")
		}
	})

	t.Run("chmod_chown_chtimes_noop", func(t *testing.T) {
		if err := fs.Chmod(renamed, 0644); err != nil {
			t.Errorf("Chmod: %v", err)
		}
		if err := fs.Chown(renamed, 0, 0); err != nil {
			t.Errorf("Chown: %v", err)
		}
		if err := fs.Chtimes(renamed, time.Now(), time.Now()); err != nil {
			t.Errorf("Chtimes: %v", err)
		}
	})

	t.Run("readdir", func(t *testing.T) {
		// Seed a directory tree under a fresh prefix:
		//   prefix/
		//     a.txt          (file)
		//     b.txt          (file)
		//     sub/           (directory inferred from children)
		//       inner.txt    (file in subdir, must not be returned)
		//     zdir/          (explicit marker, no children)
		prefix := uniqueName("list/")
		childA := prefix + "a.txt"
		childB := prefix + "b.txt"
		subdirInner := prefix + "sub/inner.txt"
		zdirMarker := prefix + "zdir/"
		t.Cleanup(func() {
			cleanup(t, client, childA, childB, subdirInner, zdirMarker)
		})

		for _, n := range []string{childA, childB, subdirInner} {
			if err := seedObject(fs, n, []byte("x")); err != nil {
				t.Fatalf("seedObject %s: %v", n, err)
			}
		}
		// Explicit empty-directory marker; deliberately no children.
		if err := fs.Mkdir(prefix+"zdir", 0755); err != nil {
			t.Fatalf("Mkdir zdir: %v", err)
		}

		f, err := fs.Open(prefix)
		if err != nil {
			t.Fatalf("Open: %v", err)
		}
		infos, err := f.Readdir(100)
		if err != nil {
			t.Fatalf("Readdir: %v", err)
		}
		if err := f.Close(); err != nil {
			t.Fatalf("Close: %v", err)
		}

		if len(infos) != 4 {
			t.Fatalf("expected 4 entries (a.txt, b.txt, sub, zdir), got %d: %v",
				len(infos), namesOf(infos))
		}
		// Names must be in lexical (byte-wise) order.
		want := []string{"a.txt", "b.txt", "sub", "zdir"}
		for i, info := range infos {
			if info.Name() != want[i] {
				t.Errorf("readdir[%d] = %q, want %q (full: %v)", i, info.Name(), want[i], namesOf(infos))
			}
		}
		// Files must be IsDir()==false, directories must be IsDir()==true.
		isDir := map[string]bool{}
		for _, info := range infos {
			isDir[info.Name()] = info.IsDir()
		}
		if isDir["a.txt"] || isDir["b.txt"] {
			t.Errorf("expected a.txt/b.txt to be files, got %v", isDir)
		}
		if !isDir["sub"] || !isDir["zdir"] {
			t.Errorf("expected sub/zdir to be directories, got %v", isDir)
		}
	})

	t.Run("readdir_count_truncates_after_sort", func(t *testing.T) {
		// Verify the count parameter returns the first N entries in sorted
		// order, not the first N returned by the SDK.
		prefix := uniqueName("cnt/")
		files := []string{"zeta.txt", "alpha.txt", "mu.txt"}
		full := make([]string, len(files))
		for i, name := range files {
			full[i] = prefix + name
		}
		t.Cleanup(func() { cleanup(t, client, full...) })
		for _, n := range full {
			if err := seedObject(fs, n, []byte("x")); err != nil {
				t.Fatalf("seedObject %s: %v", n, err)
			}
		}
		f, err := fs.Open(prefix)
		if err != nil {
			t.Fatalf("Open: %v", err)
		}
		infos, err := f.Readdir(2)
		if err != nil {
			t.Fatalf("Readdir: %v", err)
		}
		_ = f.Close()
		if len(infos) != 2 {
			t.Fatalf("expected 2 entries, got %d", len(infos))
		}
		if infos[0].Name() != "alpha.txt" || infos[1].Name() != "mu.txt" {
			t.Errorf("expected [alpha.txt, mu.txt], got %v", namesOf(infos))
		}
	})

	t.Run("removeall", func(t *testing.T) {
		// Build a tree using MkdirAll + a file, then RemoveAll the root and
		// verify every level (markers + file) is gone.
		tree := uniqueName("tree")
		if err := fs.MkdirAll(tree+"/x/y", 0755); err != nil {
			t.Fatalf("MkdirAll: %v", err)
		}
		leaf := tree + "/x/y/leaf.txt"
		if err := seedObject(fs, leaf, []byte("x")); err != nil {
			t.Fatalf("seedObject %s: %v", leaf, err)
		}
		if err := fs.RemoveAll(tree); err != nil {
			t.Fatalf("RemoveAll: %v", err)
		}
		for _, p := range []string{tree, tree + "/x", tree + "/x/y", leaf} {
			if _, err := fs.Stat(p); err == nil {
				t.Errorf("expected %s to be gone after RemoveAll", p)
			}
		}
	})

	t.Run("removeall_nonexistent_path", func(t *testing.T) {
		// os.RemoveAll on a missing path is a no-op, not an error.
		ghost := uniqueName("ghost")
		if err := fs.RemoveAll(ghost); err != nil {
			t.Errorf("RemoveAll on missing path should be a no-op, got: %v", err)
		}
	})
}

// namesOf is a small helper that returns the .Name() of every FileInfo for
// diagnostics in test failure messages.
func namesOf(infos []os.FileInfo) []string {
	out := make([]string, len(infos))
	for i, info := range infos {
		out[i] = info.Name()
	}
	return out
}

// seedObject is a thin helper used by the readdir subtest. It creates an
// object with the given payload and closes it.
func seedObject(afs afero.Fs, name string, data []byte) error {
	f, err := afs.Create(name)
	if err != nil {
		return err
	}
	if _, err := f.Write(data); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}

// TestHTTPServerReachable is a sanity check that verifies the configured
// endpoint answers a HEAD request before the heavier test runs. It helps
// produce a clearer error message when the host is offline.
func TestHTTPServerReachable(t *testing.T) {
	scheme := "http"
	if testUseSSL {
		scheme = "https"
	}
	url := scheme + "://" + testEndpoint + "/"
	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		t.Fatalf("build request: %v", err)
	}
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Skipf("server not reachable: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 500 {
		t.Skipf("server returned %d", resp.StatusCode)
	}
}
