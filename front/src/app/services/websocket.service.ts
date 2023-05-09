import { Injectable } from '@angular/core';
import { environment } from '../../environments/environment';
import { AuthService } from './auth.service';

@Injectable({
  providedIn: 'root'
})
export class WebsocketService {

  ws: WebSocket = {} as WebSocket

  constructor(private authService: AuthService) {
  }

  conn(roomID: string) {
    this.ws = new WebSocket(`ws://localhost:${environment.port}/api/rooms/${roomID}`);

    this.ws.onopen = () => {
      this.ws.onmessage = (event) => {
        if (event.data == "ping") {
          this.ws.send("pong")
          return
        }
      }

      JSON.stringify({
        email: this.authService.user?.email,
        content: "Here's some text that the server is urgently awaiting!"
      })

      this.ws.send(this.makeWSMessage("Here's some text that the server is urgently awaiting!"));
    };
  }

  sendMsg() {
    this.ws.send(this.makeWSMessage("random text"));
  }

  makeWSMessage(content: string): string {
    return JSON.stringify({
      email: this.authService.user?.email ?? '',
      content: content
    })
  }
}
