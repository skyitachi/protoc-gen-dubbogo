package dubbo

import (
  "bytes"
  "fmt"
  pb "github.com/golang/protobuf/protoc-gen-go/descriptor"
  "github.com/golang/protobuf/protoc-gen-go/generator"
  "log"
  "text/template"
)

const (
	contextPkgPath = "context"
	dubbogoPkgPath = "github.com/apache/dubbo-go"
	errorPkgPath = "github.com/pkg/errors"
	configPkgPath = "github.com/apache/dubbo-go/config"
)

var (
	contextPkg string
)

var reservedClientName = map[string]bool{
  // TODO: do we need any in gRPC?
}

type dubbogo struct {
	gen *generator.Generator
}

func init() {
  //fmt.Println("in the dubbogo plugin inited")
  //var cc *grpc.ClientConn
  //cc.Invoke()
	generator.RegisterPlugin(new(dubbogo))
}

func (d *dubbogo) P(args ...interface{}) { d.gen.P(args...) }

func (d dubbogo) Name() string {
	return "dubbogo"
}

func (d *dubbogo) Init(g *generator.Generator) {
  d.gen = g
}
// Also record that we're using it, to guarantee the associated import.
func (d *dubbogo) objectNamed(name string) generator.Object {
  d.gen.RecordTypeUse(name)
  return d.gen.ObjectNamed(name)
}

// Given a type name defined in a .proto, return its name as we will print it.
func (d *dubbogo) typeName(str string) string {
  return d.gen.TypeName(d.objectNamed(str))
}

func (d *dubbogo) GenerateImports(file *generator.FileDescriptor) {
	//fmt.Println("not implement")
}

func (d *dubbogo) Generate(file *generator.FileDescriptor) {
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}

	contextPkg = string(d.gen.AddImport(contextPkgPath))
	d.gen.AddImport(errorPkgPath)
	d.gen.AddImport(configPkgPath)

	d.P("// Reference imports to suppress errors if they are not otherwise used.")
	d.P("var _ ", contextPkg, ".Context")
	d.P()

	// Assert version compatibility.
	d.P("// This is a compile-time assertion to ensure that this generated file")
	d.P("// is compatible with the grpc package it is being compiled against.")
	d.P()

	for i, service := range file.FileDescriptorProto.Service {
		//d.generateService(file, service, i)
		d.generateService(file, service, i)
	}
}

var deprecationComment = "// Deprecated: Do not use."

func (d *dubbogo) generateService(file *generator.FileDescriptor, service *pb.ServiceDescriptorProto, index int) {
  tplStr := GetServiceTpl()
  tpl, err := template.New(service.GetName()).Parse(tplStr)
  if err != nil {
    log.Fatal(err)
  }

  serviceStub := ServiceStub{
    ServiceName: service.GetName(),
    ClientStubName: service.GetName(),
    ServerStubName: service.GetName() + "Server",
    UnimplementedServerName: "Unimplemented" + service.GetName() + "Server",
  }
  buffer := bytes.NewBuffer([]byte{})

  origServName := service.GetName()
  fullServName := origServName
  if pkg := file.GetPackage(); pkg != "" {
   fullServName = pkg + "." + fullServName
  }
  servName := generator.CamelCase(origServName)
  //deprecated := service.GetOptions().GetDeprecated()
  clientMethods := make([]string, 0)
  serverMethods := make([]string, 0)
  serverStubMethods := make([]string, 0)
  unimplementedStubMethods := make([]string, 0)
  stubName := serviceStub.ServerStubName + "Stub"

  for _, method := range service.Method {
    client, server := d.generateClientSignature(servName, method)
    unimplemented := d.generateUnimplementedServerService(server, serviceStub.UnimplementedServerName)
    clientMethods = append(clientMethods, client)
    serverMethods = append(serverMethods, server)
    unimplementedStubMethods = append(unimplementedStubMethods, unimplemented)

    methodName := method.GetName()
    returnType := d.typeName(method.GetOutputType())
    reqType := d.typeName(method.GetInputType())
    stubMethod := d.generateServerStubMethod(stubName, methodName, contextPkg, returnType, reqType)
    serverStubMethods = append(serverStubMethods, stubMethod)
  }
  serverMethods = append(serverMethods, "Reference() string")
  serviceStub.ClientMethods = clientMethods
  serviceStub.ServerMethods = serverMethods
  serviceStub.UnimplementedStubMethods = unimplementedStubMethods
  serviceStub.ServerProxyMethods = serverStubMethods
  if err := tpl.Execute(buffer, serviceStub); err != nil {
    log.Fatal(err)
  }
  d.P(string(buffer.Bytes()))
}

func (d *dubbogo) generateClientSignature(servName string, method *pb.MethodDescriptorProto) (string, string) {
  origMethName := method.GetName()
  methName := generator.CamelCase(origMethName)
  if reservedClientName[methName] {
    methName += "_"
  }
  reqArg := ", in *" + d.typeName(method.GetInputType())
  if method.GetClientStreaming() {
    reqArg = ""
  }
  respName := "out *" + d.typeName(method.GetOutputType())
  serverRespName := " *" + d.typeName(method.GetOutputType())
  if method.GetServerStreaming() || method.GetClientStreaming() {
    respName = servName + "_" + generator.CamelCase(origMethName) + "Client"
  }
  client := fmt.Sprintf("%s func (ctx %s.Context%s, %s) error", methName, contextPkg, reqArg, respName)
  server := fmt.Sprintf("%s(ctx %s.Context%s) (%s, error)", methName, contextPkg, reqArg, serverRespName)
  return client , server
}

func (d *dubbogo) generateUnimplementedServerService(methodSignature string, stubName string) string {
  return fmt.Sprintf("func (*%s) %s{\nreturn nil, errors.New(\"not implemented\")\n}", stubName, methodSignature)
}

func (d *dubbogo) generateServerStubMethod(stubName string, methodName string, ctxPkg string, returnType string, req string) string {
  tpl := GetStubMethodTpl()
  return fmt.Sprintf(tpl, stubName, methodName, ctxPkg, returnType, req, methodName)
}
