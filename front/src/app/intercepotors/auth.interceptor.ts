import { Injectable } from '@angular/core';
import { HttpEvent, HttpHandler, HttpInterceptor, HttpRequest } from '@angular/common/http';
import { Observable } from 'rxjs';
import { AuthService } from '../services/auth.service';

@Injectable()
export class AuthInterceptor implements HttpInterceptor {

  constructor(private authService: AuthService) {
  }

  intercept(request: HttpRequest<unknown>, next: HttpHandler): Observable<HttpEvent<unknown>> {
    const headers = {
      'First-name' : this.authService.user?.firstName ?? '',
      'Last-name' : this.authService.user?.lastName ?? '',
      'Email' : this.authService.user?.email ?? '',
      'User-id': this.authService.user?.id ?? ''
    }

    return next.handle(request.clone({ setHeaders: headers }));
  }
}
