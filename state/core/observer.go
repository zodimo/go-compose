package core

type ReadObserver func(source StateChangeNotifier)

type ObserverManager interface {
	WithReadObserver(ReadObserver, func())
	NotifyRead(source StateChangeNotifier)
}
