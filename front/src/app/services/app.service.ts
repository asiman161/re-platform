import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class AppService {

  constructor(private http: HttpClient) { }

  ping(): Observable<string> {
    return this.http.get(`http://${window.location.hostname}:${environment.port}/api/ping`, { responseType: 'text' })
  }
}
