import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable, Subject } from 'rxjs';
import { MediaConnection, Peer, PeerJSOption } from 'peerjs';
import { v4 as uuidv4 } from 'uuid';
import { environment } from '../../environments/environment';
import { HttpClient } from '@angular/common/http';

export interface Room {
  id: number
  name: string
  author: string
  is_open: boolean
}

@Injectable({
  providedIn: 'root'
})
export class RoomService {

  constructor(private http: HttpClient) {
  }

  public createRoom(name: string): Observable<Room> {
    return this.http.post(`http://${environment.host}:${environment.port}/api/rooms`, { name }) as Observable<Room>
  }

  public closeRoom(id: string): Observable<string> {
    return this.http.post(`http://${environment.host}:${environment.port}/api/rooms/${id}/close`, {}, { responseType: 'text' })
  }

  public getRooms(): Observable<Room[]> {
    return this.http.get(`http://${environment.host}:${environment.port}/api/rooms`) as Observable<Room[]>
  }
  public getRoom(id: string): Observable<Room> {
    return this.http.get(`http://${environment.host}:${environment.port}/api/rooms/${id}`) as Observable<Room>
  }
}
