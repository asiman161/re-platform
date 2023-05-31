import { Component, OnInit, ViewEncapsulation } from '@angular/core';
import { AuthService } from './services/auth.service';
import { SocialUser } from '@abacritt/angularx-social-login';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class AppComponent implements OnInit {

  constructor(private authService: AuthService) {}

  user: SocialUser = {} as SocialUser

  ngOnInit(): void {
    this.authService.checkAuthState()
    this.user = this.authService.user ?? {} as SocialUser
  }
}
