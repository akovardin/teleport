import {Component} from '@angular/core';
import {AuthService} from '../services/auth.service';
import {Router} from '@angular/router';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent {
  message: string;
  loading: boolean;

  constructor(public auth: AuthService, private router: Router) {
    this.loading = false;
    this.message = '';
  }

  login(email: string, password: string): boolean {
    this.loading = true;
    this.message = '';
    this.auth.login(email, password).subscribe(data => {
      this.loading = false;
      this.router.navigateByUrl('/');
    }, err => {
      if (err.error.message) {
        this.loading = false;
        this.message = err.error.message.replace(/;/g, '. ');
        setTimeout(() => {
          this.message = '';
        }, 2500);
      }
    });

    return false;
  }

  register(email: string, password: string): boolean {
    this.loading = true;
    this.message = '';
    this.auth.register(email, password).subscribe(
      data => {
        this.loading = false;
        this.router.navigateByUrl('/');
      }, err => {
        this.loading = false;
        this.message = err.error.message.replace(/;/g, '. ');
        setTimeout(() => {
          this.message = '';
        }, 2500);
      }
    );

    return false;
  }
}
