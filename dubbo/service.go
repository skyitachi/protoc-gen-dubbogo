package dubbo

func GetServiceTpl() string {
  return `
type {{.ClientStubName}} struct {
  {{range .ClientMethods}}
    {{println .}}
  {{end}}
}

func ({{.ClientStubName}}) Reference() string {
  return "{{.ServiceName}}"
}

func New{{.ClientStubName}}() *{{.ClientStubName}} {
  client := new({{.ClientStubName}})
  config.SetConsumerService(client)
  return client
}

type {{.ServerStubName}} interface {
  {{range .ServerMethods}}
    {{println .}}
  {{end}}
}

type Unimplemented{{.ServerStubName}} struct {
}

{{range .UnimplementedStubMethods }}
{{println .}}
{{end}}

func (*Unimplemented{{.ServerStubName}}) Reference() string {
  return "{{.ServiceName}}"
}

type {{.ServerStubName}}Stub struct {
  stub {{.ServerStubName}}
}

{{range .ServerProxyMethods}}
{{println .}}
{{end}}

func (*{{.ServerStubName}}Stub) Reference() string {
  return "{{.ServiceName}}"
}

func new{{.ServerStubName}}Stub(service {{.ServerStubName}}) (*{{.ServerStubName}}Stub) {
  return &{{.ServerStubName}}Stub{stub: service}
}

func RegisterProvider(service {{.ServerStubName}}) {
  stub := new{{.ServerStubName}}Stub(service)
  config.SetProviderService(stub)
}
`
}

func GetStubMethodTpl() string {
  return `func (s *%s) %s(ctx %s.Context, data []byte) (*%s, error) {
  req := &%s{}
  buf := proto.Buffer{}
  buf.SetBuf(data)
  if err := buf.Unmarshal(req); err != nil {
    return nil, err
  }
  return s.stub.%s(ctx, req)
}`
}