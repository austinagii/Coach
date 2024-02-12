<script>
  let assistantResponse = "Hey! What's your goal?";
  let userResponse = "";
  
  const goalCompletionUrl = "http://localhost:8080/goal-completion";
  async function fetchAssistantResponse(event) {
    if(event != null) {
      event.preventDefault();
    }
    let request = {
      user_message: userResponse
    }
    const response = await fetch(goalCompletionUrl, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'  
      },
      body: JSON.stringify(request)
    });

    if (response.ok) {
      const responseData = await response.json();
      assistantResponse = responseData['assistant_message']
      userResponse = ""
    }
  }
</script>

<div class="container">
  <div class="spacer-1"></div>
  <div class="assistant-response-wrapper">
    <p id="assistant-response">{assistantResponse}</p>
  </div>
  <div class="spacer-2"></div>
  <form id="user-response">
    <input type="text" placeholder="Type your response here!" bind:value={userResponse}>
    <button on:click={fetchAssistantResponse}>-></button>
  </form>
</div>

<style>
  .container {
    display: flex;
    flex-direction: column;
    align-items: center;
    width: 100vw;
    height: 100vh;
    padding: 0px;
    margin: 0px;
  }

  .container p {
    display: block;
    width: fit-content;
  }
  
  .assistant-response-wrapper {
    width: 60%;
    height: 25%;
    display: flex;
    align-items: center;
    justify-content: center;  
    padding-left: 3%;
    padding-right: 3%; 
  }

  #assistant-response {
    font-size: 24px;
    font-family: "Courier New", Times, Serif;
  }

  #user-response {
    width: 60%;
    display: flex;
    justify-content: center;
  }

  #user-response input {
    font-size: 20px;
    font-family: "Courier New", Times, Serif;
    border: none;
    border-bottom: 1px solid black;
    width: 80%;
    text-align: center;
  }

  #user-response input:focus {
    outline: none;
  }

  .spacer-1 {
    height: 20vh
  }

  .spacer-2 {
    height: 5vh
  }
</style>
