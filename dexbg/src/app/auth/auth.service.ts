import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class AuthService {

  isSignedIn: boolean = !!sessionStorage.getItem("id");
  username: string | null =  sessionStorage.getItem("username");
  dicebearFemale = `https://avatars.dicebear.com/v2/avataaars/${this.username}.svg?top%5B%5D=longHair&top%5B%5D=hat&topChance=100&accessoriesChance=40&facialHairChance=0`;
  dicebearMale = `https://avatars.dicebear.com/v2/avataaars/${this.username}.svg?mode=exclude&top%5B%5D=longHair&topChance=80&hatColor%5B%5D=pastel&hatColor%5B%5D=pink&hatColor%5B%5D=red&hairColor%5B%5D=red&hairColor%5B%5D=pastel&accessoriesChance=10&facialHairChance=80&facialHairColor%5B%5D=red&facialHairColor%5B%5D=pastel&clothesColor%5B%5D=red&clothesColor%5B%5D=pink&eyes%5B%5D=hearts`;
  noGender = `https://api.dicebear.com/5.x/thumbs/svg?seed=${this.username}`;
  profilePicture: string | null =  sessionStorage.getItem("profilePicture") || !sessionStorage.getItem("gender") ? this.noGender : parseInt(sessionStorage.getItem("gender")!) === 1 ? this.dicebearMale : this.dicebearFemale;

  constructor(private httpClient: HttpClient) { }

  signinToAccount(login: string, password: string){
    return this.httpClient.post<any>("http://localhost:3000/users/signin", {
      login: login,
      password: password
    })
  }

  joinDex(username: string, firstname: string, lastname: string, mail: string, phone: string, password: string){
    return this.httpClient.post<any>("http://localhost:3000/users/join", {
      username: username,
      password: password,
      fullname: firstname + " " + lastname,
      mail: mail,
      phone: phone, 
    })
  }

}
