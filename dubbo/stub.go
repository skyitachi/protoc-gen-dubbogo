package dubbo

type ServiceStub struct {
  ClientStubName string
  ServiceName string
  ServerStubName string
  SeverStubMethods string
  UnimplementedServerName string
  ClientMethods []string
  ServerMethods []string
  UnimplementedStubMethods []string
  ServerProxyMethods []string
}
