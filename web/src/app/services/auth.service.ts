import {Inject, Injectable} from '@angular/core';
import {HttpClient, HttpHeaders} from '@angular/common/http';
import {Observable} from 'rxjs';
import {tap} from 'rxjs/operators';

class Auth {
  email: string;
  token: string;
}

@Injectable()
export class AuthService {
  headers: HttpHeaders;

  constructor(@Inject('api') private api: string,
              private http: HttpClient) {
    this.headers = new HttpHeaders();
    this.headers = this.headers.append('Content-Type', 'application/json');
  }

  login(email: string, password: string): Observable<any> {
    const observer = this.http.post(this.api + 'users/login', JSON.stringify({
      email: email,
      password: password,
    }), {headers: this.headers});
    return observer.pipe(
      tap((data: Auth) => {
        if (data.token !== '') {
          localStorage.setItem('email', data.email);
          localStorage.setItem('token', data.token);
        }
      })
    );
  }

  register(email: string, password: string): Observable<any> {
    const observer = this.http.post(this.api + 'users/register', JSON.stringify({
      email: email,
      password: password,
    }), {headers: this.headers});
    return observer.pipe(
      tap((data: Auth) => {
          console.log(data);
          if (data.token !== '') {
            this.setParams(data.email, data.token);
          }
        },
        err => {
          console.log(err);
        }),
    );
  }

  update(password: string): Observable<any> {
    const observer = this.http.post(this.api + 'users/update', JSON.stringify({
      password: password,
    }), {headers: this.authHeaders()});
    return observer.pipe(
      tap(data => {
          console.log(data);
        },
        err => {
          console.log(err);
        }),
    );
  }

  restore(email: string): Observable<any> {
    const observer = this.http.post(this.api + 'users/restore', JSON.stringify({
      email: email,
    }), {headers: this.headers});
    return observer.pipe(
      tap(data => {
          console.log(data);
        },
        err => {
          console.log(err);
        }),
    );
  }

  logout(): any {
    localStorage.removeItem('email');
    localStorage.removeItem('token');
  }

  getEmail(): string {
    return localStorage.getItem('email');
  }

  getToken(): string {
    return localStorage.getItem('token');
  }

  setParams(email: string, token: string) {
    localStorage.setItem('email', email);
    localStorage.setItem('token', token);
  }

  isLoggedin(): boolean {
    console.log(this.getEmail(), this.getToken())
    return this.getEmail() !== null && this.getToken() !== null;
  }

  authHeaders(): HttpHeaders {
    let headers = new HttpHeaders();
    headers = headers.append('Content-Type', 'application/json');
    return headers.append('Authorization', 'Bearer ' + this.getToken());
  }
}
