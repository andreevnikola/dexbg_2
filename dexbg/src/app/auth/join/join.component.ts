import { HttpErrorResponse } from '@angular/common/http';
import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { faExclamationCircle } from '@fortawesome/free-solid-svg-icons';
import { AuthService } from '../auth.service';

@Component({
  selector: 'app-join',
  templateUrl: './join.component.html',
  styleUrls: ['./join.component.scss']
})
export class JoinComponent {

  passwordIsShwon: boolean = false;
  phoneStarting: any = "+359";
  loading: boolean = false;
  error: boolean | number = false;
  hadAnError: boolean = false;
  nottaBenneIcon = faExclamationCircle;

  constructor( private authSerrvice: AuthService, private router: Router ) { }

  handleSigninUser(username: string, firstname: string, lastname: string, mail: string, phone: string, password: string){
    this.loading = true;
    this.error = false;
    this.authSerrvice.joinDex(username, firstname, lastname, mail, this.phoneStarting + phone, password).subscribe({
      next: (data: any) => {
        this.error = false;
        localStorage.setItem("key", data.key);
        sessionStorage.setItem("id", data.id);
        sessionStorage.setItem("username", username);
        sessionStorage.setItem("mail", mail);
        sessionStorage.setItem("phone", this.phoneStarting + phone);
        sessionStorage.setItem("fullname", firstname + lastname);
        sessionStorage.setItem("gender", data.gender);
        this.loading = false;
        this.authSerrvice.isSignedIn = true;
        const dicebearFemale = `https://avatars.dicebear.com/v2/avataaars/${username}.svg?top%5B%5D=longHair&top%5B%5D=hat&topChance=100&accessoriesChance=40&facialHairChance=0`;
        const dicebearMale = `https://avatars.dicebear.com/v2/avataaars/${username}.svg?mode=exclude&top%5B%5D=longHair&topChance=80&hatColor%5B%5D=pastel&hatColor%5B%5D=pink&hatColor%5B%5D=red&hairColor%5B%5D=red&hairColor%5B%5D=pastel&accessoriesChance=10&facialHairChance=80&facialHairColor%5B%5D=red&facialHairColor%5B%5D=pastel&clothesColor%5B%5D=red&clothesColor%5B%5D=pink&eyes%5B%5D=hearts`;
        const noGender = `https://api.dicebear.com/5.x/thumbs/svg?seed=${username}`;
        this.authSerrvice.profilePicture = !data.gender ? noGender : data.gender === 1 ? dicebearMale : dicebearFemale;
        this.authSerrvice.username = username;
        this.router.navigate( ["/profile/" + data.id] );
      },
      error: (err: HttpErrorResponse) => {
        this.loading = false;
        if((err as any).handled){ return }
        this.error = err.status;
      }
    });
  }

}
