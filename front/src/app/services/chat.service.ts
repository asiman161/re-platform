import { Injectable } from '@angular/core';
import { environment } from '../../environments/environment';
import { Observable } from 'rxjs';
import { HttpClient } from '@angular/common/http';

export interface ChatMessage {
  id: number
  room_id: string
  content: string
  author: string
  created_at: string
}

@Injectable({
  providedIn: 'root'
})
export class ChatService {

  constructor(private http: HttpClient) { }

  getMessages(roomID: string): Observable<ChatMessage[]> {
    return this.http.get(`http://${window.location.hostname}:${environment.port}/api/rooms/${roomID}/chat`) as Observable<ChatMessage[]>
  }

  sendMessage(roomID: string, msg: string): Observable<ChatMessage[]> {
    return this.http.post(`http://${window.location.hostname}:${environment.port}/api/rooms/${roomID}/chat`, { content: msg }) as Observable<ChatMessage[]>
  }
}
