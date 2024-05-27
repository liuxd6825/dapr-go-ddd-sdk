package class

type Access string

/*
	-	private 私有
	#	protected 受保护
	~	package private 包内可见
	+	public 公有
*/

const (
	AccessPublic         Access = "public"
	AccessPrivate        Access = "private"
	AccessProtected      Access = "protected"
	AccessPackagePrivate Access = "package private"
)
