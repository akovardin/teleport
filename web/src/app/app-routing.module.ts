import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import {LoggedGuard} from "./logged.guard";

import {LoginComponent} from './login/login.component';
import {IntegrationsComponent} from "./integrations/integrations.component";
import {UsersComponent} from "./users/users.component";

const routes: Routes = [
  {path: '', redirectTo: 'integrations', pathMatch: 'full'},
  {path: 'login', component: LoginComponent},
  {path: 'integrations', component: IntegrationsComponent, canActivate: [LoggedGuard]},
  {path: 'users', component: UsersComponent, canActivate: [LoggedGuard]},
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
