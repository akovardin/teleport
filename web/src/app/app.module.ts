import {BrowserModule} from '@angular/platform-browser';
import {NgModule} from '@angular/core';
import {HttpClientModule} from '@angular/common/http';
import {FormsModule} from '@angular/forms';

import {AppRoutingModule} from './app-routing.module';

import {environment} from '../environments/environment';

import {LoggedGuard} from './logged.guard';
import {LoginComponent} from './login/login.component';
import {AuthService} from './services/auth.service';
import {AppComponent} from './app.component';

const server = environment.server;

@NgModule({
  declarations: [
    AppComponent,
    LoginComponent,
  ],
  imports: [
    BrowserModule,
    HttpClientModule,
    FormsModule,
    AppRoutingModule
  ],
  providers: [
    AuthService,
    LoggedGuard,
    {provide: 'server', useValue: server},
  ],
  bootstrap: [AppComponent]
})
export class AppModule {
}
