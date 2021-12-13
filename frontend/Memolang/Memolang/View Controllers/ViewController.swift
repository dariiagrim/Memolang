//
//  ViewController.swift
//  Memolang
//
//  Created by Dariia Hrymalska on 20.10.2021.
//

import UIKit

class ViewController: UIViewController {
    @IBOutlet weak var signUpButton: UIButton!
    @IBOutlet weak var loginButton: UIButton!
    @IBOutlet var logoBlocks: [UIView]!
    @IBOutlet var superview: UIView!
    
    override func viewDidLoad() {
        super.viewDidLoad()
        var constraints = [NSLayoutConstraint]()
        logoBlocks.forEach { block in
            block.layer.cornerRadius = 45
            block.translatesAutoresizingMaskIntoConstraints = false
            constraints.append(block.centerXAnchor.constraint(equalTo: superview.centerXAnchor))
            constraints.append(block.centerYAnchor.constraint(equalTo: superview.centerYAnchor))
        }
        NSLayoutConstraint.activate(constraints)
    }
}

