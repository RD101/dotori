package main

import (
	"errors"
	"os"
	"regexp"
	"strings"
)

var (
	regexRFC3339Time      = regexp.MustCompile(`^\d{4}-(0?[1-9]|1[012])-(0?[1-9]|[12][0-9]|3[01])T([0-1][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9][-+]\d{2}:\d{2}$`)               //2019-09-09T02:46:52+09:00
	regexPath             = regexp.MustCompile(`^/[/ _.가-힣\w]+$`)                                                                                                           // 일반경로 /LIBRARY_3D/asset/
	regexSingleQuotesPath = regexp.MustCompile(`'/[/ _.가-힣\w]+'`)                                                                                                           // 작은 따옴표로 구성된 경로 '/LIBRARY_3D/asset/'
	regexDoubleQuotesPath = regexp.MustCompile(`"/[/ _.가-힣\w]+"`)                                                                                                           // 큰 따옴표로 구성된 경로 "/LIBRARY_3D/asset/"
	regexIPv4             = regexp.MustCompile(`^([01]?\d?\d|2[0-4]\d|25[0-5]).([01]?\d?\d|2[0-4]\d|25[0-5]).([01]?\d?\d|2[0-4]\d|25[0-5]).([01]?\d?\d|2[0-4]\d|25[0-5])$`) // 0.0.0.0 ~ 255.255.255.255
	regexLower            = regexp.MustCompile(`[a-z0-9]+$`)                                                                                                                // Itemtype, Dbname (maya, nuke, fusion360..)
	regexObjectID         = regexp.MustCompile(`^[a-z0-9]*$`)                                                                                                               // "54759eb3c090d83494e2d804"
	regexMap              = regexp.MustCompile(`^([a-zA-Z0-9]+:[a-zA-Z0-9.-_]+)(,?([a-zA-Z0-9]+:[a-zA-Z0-9.-_]+))*$`)                                                       // key:value,key:value
	regexTag              = regexp.MustCompile(`^[가-힣a-zA-Z0-9]+$`)                                                                                                         // 태그, tag, tag1
	regexPermission       = regexp.MustCompile(`^[0][0-7][0-7][0-7]$`)                                                                                                      //0775, 0440 (권한은 0000~7777까지 가능하지만 보안상 0000~0777까지만 허용한다)
	regexSplitBySign      = regexp.MustCompile(`[,/ _]+`)
	regexTitle            = regexp.MustCompile(`^[가-힣\w\s]+$`)
)

//FileExists 함수는 해당 파일이 존재하는지 체크한다.
func FileExists(path string) error {
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("파일명을 정확히 입력해주세요")
		}
		return err
	}
	if stat.IsDir() {
		return errors.New("파일명까지 입력해주세요")
	}
	return nil
}

func str2bool(str string) bool {
	return strings.ToLower(str) == "true"
}
