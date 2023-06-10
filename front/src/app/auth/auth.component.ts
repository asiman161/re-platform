import { Component } from '@angular/core';
import { AuthService } from '../services/auth.service';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-auth',
  templateUrl: './auth.component.html',
  styleUrls: ['./auth.component.scss']
})
export class AuthComponent {
  constructor(private authService: AuthService, private activatedRoute: ActivatedRoute) {
    this.activatedRoute.queryParams.subscribe({
      next: data => {
        if (!!data['email'] && !this.authService.user) {
          this.authService.setUserEmail(data['email'])
        }
      }
    })
  }

  signOut(): void {
    this.authService.signOut();
  }
}
