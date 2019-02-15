export class Integration {
  id?: number;
  title?: string;
  token?: string;
  channel?: string;
  secret?: string;

  proxyAddress?: string;
  proxyUser?: string;
  proxyPass?: string;


  constructor(id?: number, title?: string, token?: string, channel?: string, secret?: string, proxyAddress?: string, proxyUser?: string, proxyPass?: string) {
    this.id = id;
    this.title = title;
    this.token = token;
    this.channel = channel;
    this.secret = secret;
    this.proxyAddress = proxyAddress;
    this.proxyUser = proxyUser;
    this.proxyPass = proxyPass;
  }
}
