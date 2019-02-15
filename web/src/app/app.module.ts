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
import {IntegrationsComponent} from './integrations/integrations.component';
import {UsersComponent} from './users/users.component';
import {IntegrationsService} from "./integrations/integrations.service";

const api = environment.api;
const domain = environment.domain;

@NgModule({
  declarations: [
    AppComponent,
    LoginComponent,
    IntegrationsComponent,
    UsersComponent,
  ],
  imports: [
    BrowserModule,
    HttpClientModule,
    FormsModule,
    AppRoutingModule
  ],
  providers: [
    IntegrationsService,
    AuthService,
    LoggedGuard,
    {provide: 'api', useValue: api},
    {provide: 'domain', useValue: domain},
  ],
  bootstrap: [AppComponent]
})
export class AppModule {
}
