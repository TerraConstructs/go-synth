import { App, HttpBackend } from "cdktf";
import { aws } from "terraconstructs";

const outdir = "cdktf.out";
const app = new App({
  outdir,
});
const stack = new aws.AwsStack(app, "network-stack", {
  gridUUID: "12345678-1234",
  environmentName: "test",
  providerConfig: {
    region: "us-west-2",
  },
});

new aws.network.SimpleIPv4Vpc(stack, "network", {
  internalDomain: "example.com",
  ipv4CidrBlock: "10.0.0.0/16",
});

new HttpBackend(stack, {
  address: "http://localhost:1234",
});

app.synth();
