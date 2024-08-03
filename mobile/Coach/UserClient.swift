import Foundation

let USER_API_URL = "http://localhost:8080/users"

/// Represents a user.
struct User: Codable {
    let id: String
    let name: String
}

// Represents a request to create a new user.
struct NewUserRequest: Codable {
    let name: String
}

/// Creates a new user with the given name.
///
/// - Parameter name: The user's desired name
/// - Returns: The newly created user
func createUser(_ name: String) async -> User? {
    let request = NewUserRequest(name: name)
    guard let url = URL(string: USER_API_URL) else {
        log.error("Could not create URL for user API")
        return nil
    }
    
    guard let user: User = await performHttpPostRequest(url: url, requestBody: request) else {
        log.warning("Failed to create user")
        return nil
    }
    
    log.info("User created successfully")
    return user
}
