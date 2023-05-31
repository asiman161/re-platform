import { Component, OnInit } from '@angular/core';
import { AuthService } from '../services/auth.service';
import { SocialUser } from '@abacritt/angularx-social-login';

@Component({
  selector: 'app-navbar',
  templateUrl: './navbar.component.html',
  styleUrls: ['./navbar.component.scss']
})
export class NavbarComponent implements OnInit {
  constructor(private authService: AuthService) {
  }

  user: SocialUser = {} as SocialUser

  ngOnInit(): void {
    this.user = this.authService.user ?? {} as SocialUser
  }
}
