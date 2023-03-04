import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { SigninComponent } from './signin/signin.component';
import { JoinComponent } from './join/join.component';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { FormsModule } from '@angular/forms';
import { AuthRoutingModule } from './auth-routing.module';
import { SharedModule } from '../shared/shared.module';



@NgModule({
  declarations: [
    SigninComponent,
    JoinComponent
  ],
  imports: [
    CommonModule,
    FontAwesomeModule,
    FormsModule,
    AuthRoutingModule,
    SharedModule
  ],
  exports: [
    SigninComponent,
    JoinComponent
  ]
})
export class AuthModule { }
