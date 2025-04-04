package transform

var tsc *Tsc

var babel *Babel

// TransformFromTypeScript
//
//	@Description: 将typescript代码转换为js
//	@param tsCode
//	@return string
//	@return error
func TransformFromTypeScript(tsCode string, fileName string) ([]byte, error) {
	if tsc == nil {
		tsc = NewTsc()
	}
	es5Code, err := tsc.TransformEs5(tsCode, fileName)
	return es5Code, err
}

// TransformFromEs6
//
//	@Description:
//	@param es6Code
//	@param fileName
//	@return []byte
//	@return error
func TransformFromEs6(es6Code string, fileName string) ([]byte, error) {
	if babel == nil {
		var err error
		if babel, err = NewBabel(); err != nil {
			return nil, err
		}
	}
	codeBytes, resErr := babel.Transform(es6Code, fileName, false, nil)
	return codeBytes, resErr
}
