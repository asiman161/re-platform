import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';
import { HttpClient } from '@angular/common/http';
import { makePathPrefix } from './utils';

export interface Room {
  id: number
  name: string
  author: string
  is_open: boolean
  created_at: string
  updated_at: string
}

export interface Quiz {
  id: number
  room_id: string
  name: string
  author: string
  content: string
  variants: Variant[]
  answers: Answer[]
  is_open: boolean
  created_at?: string
  updated_at?: string
}

export interface Variant {
  id: number
  value: string
}

export interface Answer {
  variant_id: number
  author: string
}

@Injectable({
  providedIn: 'root'
})
export class RoomService {

  constructor(private http: HttpClient) {
  }

  public createQuiz(quiz: Partial<Quiz>): Observable<Quiz> {
    const path = `${environment.schema}${environment.host}:${environment.port}/api/rooms/${quiz.room_id}/quizzes`
    return this.http.post(path, quiz) as Observable<Quiz>
  }

  public getQuizzes(roomID: string): Observable<Quiz[]> {
    const path = `${makePathPrefix()}/api/rooms/${roomID}/quizzes`
    return this.http.get(path) as Observable<Quiz[]>
  }

  public answerQuiz(quizID: number, variant: Variant): Observable<string> {
    const path = `${makePathPrefix()}/api/rooms/-/quizzes/${quizID}/answer`
    return this.http.post(path, { 'variant_id': variant.id }, { responseType: 'text' }) as Observable<string>
  }

  public closeQuiz(quizID: number): Observable<string> {
    const path = `${makePathPrefix()}/api/rooms/-/quizzes/${quizID}/close`
    return this.http.post(path, {}, { responseType: 'text' }) as Observable<string>
  }

  public createRoom(name: string): Observable<Room> {
    return this.http.post(`${environment.schema}${environment.host}:${environment.port}/api/rooms`, { name }) as Observable<Room>
  }

  public closeRoom(id: string): Observable<string> {
    return this.http.post(`${environment.schema}${environment.host}:${environment.port}/api/rooms/${id}/close`, {}, { responseType: 'text' })
  }

  public getRooms(): Observable<Room[]> {
    return this.http.get(`${environment.schema}${environment.host}:${environment.port}/api/rooms`) as Observable<Room[]>
  }

  public getRoom(id: string): Observable<Room> {
    return this.http.get(`${environment.schema}${environment.host}:${environment.port}/api/rooms/${id}`) as Observable<Room>
  }
}
