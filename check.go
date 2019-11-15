package main

import "regexp"

var (
	regexRFC3339Time = regexp.MustCompile(`^\d{4}-(0?[1-9]|1[012])-(0?[1-9]|[12][0-9]|3[01])T([0-1][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9][-+]\d{2}:\d{2}$`)               //2019-09-09T02:46:52+09:00
	regexPath        = regexp.MustCompile(`^/[/_.\w]+$`)                                                                                                               // /LIBRARY_3D/asset/
	regexIPv4        = regexp.MustCompile(`^([01]?\d?\d|2[0-4]\d|25[0-5]).([01]?\d?\d|2[0-4]\d|25[0-5]).([01]?\d?\d|2[0-4]\d|25[0-5]).([01]?\d?\d|2[0-4]\d|25[0-5])$`) // 0.0.0.0 ~ 255.255.255.255
	regexLower       = regexp.MustCompile(`[a-z]+$`)                                                                                                                   // Itemtype, Dbname (maya, nuke, dotori..)
)
