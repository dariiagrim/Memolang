//
//  SignUpViewController.swift
//  Memolang
//
//  Created by Dariia Hrymalska on 07.12.2021.
//

import UIKit
import FirebaseAuth
import Alamofire

class SignUpViewController: UIViewController {

    @IBOutlet weak var firstNameTextField: UITextField!
    @IBOutlet weak var lastNameTextField: UITextField!
    @IBOutlet weak var emailTextField: UITextField!
    @IBOutlet weak var passwordTextField: UITextField!
    @IBOutlet weak var signUpButton: UIButton!
    @IBOutlet weak var errorLabel: UILabel!
    
    
    override func viewDidLoad() {
        super.viewDidLoad()

        // Do any additional setup after loading the view.
    }
    
    @IBAction func signUpTapped(_ sender: Any) {
        errorLabel.text = ""
        let email = emailTextField.text?.trimmingCharacters(in: .whitespacesAndNewlines) ?? ""
        let password = passwordTextField.text?.trimmingCharacters(in: .whitespacesAndNewlines) ?? ""
        let firstName = firstNameTextField.text?.trimmingCharacters(in: .whitespacesAndNewlines) ?? ""
        let lastName = lastNameTextField.text?.trimmingCharacters(in: .whitespacesAndNewlines) ?? ""
        
        Auth.auth().createUser(withEmail: email, password: password) { [weak self] result, error in
            if let error = error {
                self?.errorLabel.text = error.localizedDescription
                return
            }
            let user = User(firstName: firstName, lastName: lastName)
            Auth.auth().currentUser?.getIDToken(completion: { token, error in
                if let token = token {
                    print(token)
                    AF.request("https://google.com", method: .post, parameters: ["user": user], encoding: JSONEncoding(), headers: [HTTPHeader(name: "Authorization", value: "Bearer \(token)")], interceptor: nil, requestModifier: nil).responseData { response in
                        if let data = response.data {
                            let decoder = JSONDecoder()
                            let user = try? decoder.decode(User.self, from: data)
                        }
                    }
                } else {
                    self?.errorLabel.text = error?.localizedDescription
                }
            })
        }
    }
}

struct User: Codable {
    var firstName: String
    var lastName: String
}
