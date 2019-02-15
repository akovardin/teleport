import {Component, OnInit} from '@angular/core';
import {Router} from '@angular/router';
import {AuthService} from './services/auth.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit {
  activeLinkIndex = -1;

  constructor(private router: Router, public auth: AuthService) {
  }

  logout(): boolean {
    this.auth.logout();
    this.router.navigate(['/login']);
    return false;
  }

  ngOnInit() {
    this.router.events.subscribe(res => {
      if (this.router.url === '/integrations') {
        this.activeLinkIndex = 1;
      } else if (this.router.url === '/users') {
        this.activeLinkIndex = 2;
      }
    });
  }
}
