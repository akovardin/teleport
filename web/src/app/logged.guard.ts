import {Injectable} from '@angular/core';
import {CanActivate, ActivatedRouteSnapshot, RouterStateSnapshot, Router} from '@angular/router';
import {Observable} from 'rxjs';
import {AuthService} from './services/auth.service';

@Injectable({
  providedIn: 'root'
})
export class LoggedGuard implements CanActivate {
  constructor(private auth: AuthService, private router: Router) {
  }

  canActivate(
    next: ActivatedRouteSnapshot,
    state: RouterStateSnapshot): Observable<boolean> | Promise<boolean> | boolean {
    const isLoggedIn = this.auth.isLoggedin();
    console.log('canActivate', isLoggedIn);

    if (!isLoggedIn) {
      this.router.navigate(['/login'], {
        queryParams: {
          return: state.url
        }
      });
    }

    return isLoggedIn;
  }
}
