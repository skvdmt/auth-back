package model

// Errors Глобальный канал ошибок.
var Errors chan error

// ExternalErrors Глобальный канал внешних ошибок не влияющих
// на работу приложения, но записывающихся в журнал.
var ExternalErrors chan error
