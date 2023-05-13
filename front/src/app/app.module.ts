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
import { StreamComponent } from './room/stream/stream.component';
import { CallInfoDialogComponents } from './room/stream/dialog/callinfo-dialog.component';
import { MatFormFieldModule } from '@angular/material/form-field';
import { ClipboardModule } from '@angular/cdk/clipboard';
import { MatDialogModule } from '@angular/material/dialog';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatInputModule } from '@angular/material/input';
import { ChatComponent } from './room/chat/chat.component';
import { MatCardModule } from '@angular/material/card';
import { ChatService } from './services/chat.service';
import { OnlineComponent } from './room/online/online.component';

@NgModule({
  bootstrap: [AppComponent],
  declarations: [
    AppComponent,
    DashboardComponent,
    RoomComponent,
    AuthComponent,
    StreamComponent,
    CallInfoDialogComponents,
    ChatComponent,
    OnlineComponent,
  ],
    imports: [
        BrowserModule,
        AppRoutingModule,
        HttpClientModule,
        FormsModule,
        BrowserAnimationsModule,
        MatToolbarModule,
        MatIconModule,
        MatButtonModule,
        MatSnackBarModule,
        SocialLoginModule,
        GoogleSigninButtonModule,
        MatFormFieldModule,
        ClipboardModule,
        MatDialogModule,
        MatInputModule,
        MatCardModule,
        ReactiveFormsModule,
    ],
  providers: [
    AppService,
    RoomService,
    AuthGuardService,
    AuthService,
    ChatService,
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
