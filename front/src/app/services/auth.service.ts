import { Injectable } from '@angular/core';
import { SocialAuthService, SocialUser } from '@abacritt/angularx-social-login';
import { Router } from '@angular/router';

@Injectable({
  providedIn: 'root'
})
export class AuthService {

  constructor(private router: Router, private socialAuthService: SocialAuthService) { }

  private userData: SocialUser | undefined;

  get user() {
    return this.userData
  }

  checkAuthState(): boolean {
    let localUser = localStorage.getItem('user')
    if (!!localUser) {
      this.userData = JSON.parse(localUser)
      return true
    }
    let flag = false;
    this.socialAuthService.authState.subscribe({
      next: (socialUser: SocialUser) => {
        if (socialUser) {
          this.userData = socialUser
          localStorage.setItem('user', JSON.stringify(socialUser))
          flag = true;
        }
      },
      error: (error) => {
        console.log(error);
      },
      complete: () => {
        console.log('Complete!');
      },
    });
    return flag;
  }

  signOut() {
    this.userData = undefined
    localStorage.removeItem('user')
    this.router.navigateByUrl('/')
  }
}
