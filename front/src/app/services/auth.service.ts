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

  setUserEmail(email: string) {
    if (!this.userData) {
      this.userData = {} as SocialUser
    }
    this.userData.email = email
    this.saveUser(this.userData)
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
          this.saveUser(socialUser)
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

  saveUser(u: SocialUser) {
    localStorage.setItem('user', JSON.stringify(u))
  }

  signOut() {
    this.userData = undefined
    localStorage.removeItem('user')
    this.router.navigateByUrl('/')
  }
}
