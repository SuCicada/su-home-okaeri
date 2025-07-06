package logger

import "go.uber.org/zap"

var Logger, _ = zap.NewDevelopment()
var Sugar = Logger.Sugar()
var Info = Sugar.Info
var Debug = Sugar.Debug
var Warn = Sugar.Warn
var Error = Sugar.Error
