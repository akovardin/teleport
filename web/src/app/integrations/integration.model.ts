export class Integration {
  id?: number;
  title?: string;
  botKey?: string;
  channelName?: string;
  vkSecret?: string;


  constructor(id?: number, title?: string, botKey?: string, channelName?: string, vkSecret?: string) {
    this.id = id;
    this.title = title;
    this.botKey = botKey;
    this.channelName = channelName;
    this.vkSecret = vkSecret;
  }
}
