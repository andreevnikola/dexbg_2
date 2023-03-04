import { HttpErrorResponse } from '@angular/common/http';
import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { faExclamationCircle } from '@fortawesome/free-solid-svg-icons';
import { AuthService } from '../auth.service';

@Component({
  selector: 'app-signin',
  templateUrl: './signin.component.html',
  styleUrls: ['./signin.component.scss']
})
export class SigninComponent {

  passwordIsShwon: boolean = false;
  loading: boolean = false;
  error: boolean = false;
  hadAnError: boolean = false;
  nottaBenneIcon = faExclamationCircle;

  constructor( private authSerrvice: AuthService, private router: Router ) { }

  handleSigninUser(login: string, password: string){
    this.loading = true;
    this.error = false;
    this.authSerrvice.signinToAccount(login, password).subscribe({
      next: (data: any) => {
        this.error = false;
        localStorage.setItem("key", data.key);
        sessionStorage.setItem("id", data.id);
        sessionStorage.setItem("username", data.username);
        sessionStorage.setItem("mail", data.mail);
        sessionStorage.setItem("phone", data.phone);
        sessionStorage.setItem("fullname", data.fullname);
        sessionStorage.setItem("gender", data.gender);
        this.loading = false;
        this.authSerrvice.isSignedIn = true;
        const dicebearFemale = `https://avatars.dicebear.com/v2/avataaars/${data.username}.svg?top%5B%5D=longHair&top%5B%5D=hat&topChance=100&accessoriesChance=40&facialHairChance=0`;
        const dicebearMale = `https://avatars.dicebear.com/v2/avataaars/${data.username}.svg?mode=exclude&top%5B%5D=longHair&topChance=80&hatColor%5B%5D=pastel&hatColor%5B%5D=pink&hatColor%5B%5D=red&hairColor%5B%5D=red&hairColor%5B%5D=pastel&accessoriesChance=10&facialHairChance=80&facialHairColor%5B%5D=red&facialHairColor%5B%5D=pastel&clothesColor%5B%5D=red&clothesColor%5B%5D=pink&eyes%5B%5D=hearts`;
        const noGender = `https://api.dicebear.com/5.x/thumbs/svg?seed=${data.username}`;
        this.authSerrvice.profilePicture = data.profilePicture || !data.gender ? noGender : data.gender === 1 ? dicebearMale : dicebearFemale;
        this.authSerrvice.username = data.username;
        this.router.navigate( ["/profile/" + data.id] );
      },
      error: (err: HttpErrorResponse) => {
        this.loading = false;
        if((err as any).handled){ return }
        this.error = true;
      }
    });
  }
  
  isEmail(field: string){
    const reg = /^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+$/;
    return reg.test(field);
  }

}
