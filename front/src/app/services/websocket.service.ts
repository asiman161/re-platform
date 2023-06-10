import { Injectable } from '@angular/core';
import { environment } from '../../environments/environment';
import { AuthService } from './auth.service';

import { webSocket, WebSocketSubject } from 'rxjs/webSocket';

@Injectable({
  providedIn: 'root'
})
export class WebsocketService {

  // @ts-ignore
  subject: WebSocketSubject<any>

  constructor(private authService: AuthService) {
  }

  conn(roomID: string, email: string) : WebSocketSubject<any> {
    const path = `ws://${window.location.hostname}:${environment.port}/api/rooms/${roomID}/ws/${email}`
    this.subject = webSocket(path);

    return this.subject
  }



  sendMsg(msg: string) {
    this.subject.next({
      email: this.authService.user?.email ?? '',
      content: msg
    })
  }
}
