import { Component, ViewChild } from '@angular/core';
import { faBars } from '@fortawesome/free-solid-svg-icons';
import {
  trigger,
  style,
  animate,
  transition,
  // ...
} from '@angular/animations';
import { Router } from '@angular/router';
import { AuthService } from 'src/app/auth/auth.service';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.scss'],
  animations: [
    trigger('menu', [
      transition(':enter', [style({
        opacity: 0,
        height: 0
      }), animate('250ms ease-in', style({opacity: 1, height: '*'}))]),
      transition(':leave', [style({
        opacity: 1,
        'height': '*'
      }), animate('150ms ease-in', style({opacity: 0, height: 0}))]),
    ]),
    trigger('profile', [
      transition(':enter', [style({
        opacity: 0,
      }), animate('250ms ease-in', style({opacity: 1}))]),
      transition(':leave', [style({
        opacity: 1,
        scale: 1
      }), animate('150ms ease-in', style({opacity: 0, scale: 1.1}))]),
    ])
  ]
})
export class HeaderComponent {

  bars = faBars;
  mobileMenuShown: boolean = false;
  profileDataShown: boolean = false;

  @ViewChild("accSettings") accSettings: any;

  constructor( public router: Router, public authService: AuthService ) { }

  showHideMenuHandler(){
    this.mobileMenuShown = !this.mobileMenuShown;
  }

  showProfileSettings(){
    this.profileDataShown = true;
    this.mobileMenuShown = false;
  }

}
