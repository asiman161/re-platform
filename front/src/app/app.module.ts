import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { HTTP_INTERCEPTORS, HttpClientModule } from '@angular/common/http';
import { DashboardComponent } from './dashboard/dashboard.component';
import { RoomComponent } from './room/room.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { AppService } from './services/app.service';
import { MAT_SNACK_BAR_DEFAULT_OPTIONS, MatSnackBarModule } from '@angular/material/snack-bar';
import { RoomService } from './services/room.service';
import {
  GoogleLoginProvider,
  GoogleSigninButtonModule,
  SocialAuthServiceConfig,
  SocialLoginModule
} from '@abacritt/angularx-social-login';
import { AuthComponent } from './auth/auth.component';
import { AuthGuardService } from './services/auth-guard.service';
import { AuthService } from './services/auth.service';
import { AuthInterceptor } from './intercepotors/auth.interceptor';

@NgModule({
  bootstrap: [AppComponent],
  declarations: [
    AppComponent,
    DashboardComponent,
    RoomComponent,
    AuthComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    HttpClientModule,
    BrowserAnimationsModule,
    MatToolbarModule,
    MatIconModule,
    MatButtonModule,
    MatSnackBarModule,
    SocialLoginModule,
    GoogleSigninButtonModule,
  ],
  providers: [
    AppService,
    RoomService,
    AuthGuardService,
    AuthService,
    { provide: HTTP_INTERCEPTORS, useClass: AuthInterceptor, multi: true },
    {provide: MAT_SNACK_BAR_DEFAULT_OPTIONS, useValue: {duration: 2500}},
    {
      provide: 'SocialAuthServiceConfig',
      useValue: {
        autoLogin: false,
        providers: [
          {
            id: GoogleLoginProvider.PROVIDER_ID,
            provider: new GoogleLoginProvider(
              '632361311309-69vjg3si7db8618luvs71jdf7g7ugsak.apps.googleusercontent.com'
            )
          },
        ],
      } as SocialAuthServiceConfig,
    }
  ]
})
export class AppModule { }
