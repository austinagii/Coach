import Foundation

/// Performs an HTTP request to the given url with request body encoded as JSON.
///
/// - Parameters:
///   - url: The URL that the request should be sent to
///   - requestBody: A c
func performHttpPostRequest<Req: Codable, Res: Codable>(url: URL, requestBody: Req) async -> Res? {
    guard let requestBody = try? JSONEncoder().encode(requestBody) else {
        return nil
    }
    
    var request = URLRequest(url: url)
    request.httpMethod = "POST"
    request.setValue("application/json", forHTTPHeaderField: "Content-Type")
    
    do {
        let (data, response) = try await URLSession.shared.upload(for: request, from: requestBody)
        let responseStatusCode = (response as? HTTPURLResponse)?.statusCode ?? -1
        guard (responseStatusCode >= 200 && responseStatusCode <= 299) else {
            print("Request was unsuccessful, returned status code \(responseStatusCode)")
            return nil
        }
        guard let responseBody = try? JSONDecoder().decode(Res.self, from: data) else {
            return nil
        }
        return responseBody
    } catch {
        print("Error: \(error)")
        return nil
    }
}
