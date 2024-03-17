<script>
  import { onMount } from 'svelte';

  let chatId = ""
  let assistantMessage = "";
  let userMessage = "";

  async function startNewChat() {
    const response = await fetch("http://localhost:8080/chats", {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'  
      },
      body: JSON.stringify({
        task: "goal_creation"
      })
    });

    if (response.ok) {
      const responseData = await response.json();
      chatId = responseData['id']
      assistantMessage = responseData['text']
    }
  } 

  async function fetchAssistantResponse(event) {
    if(event != null) {
      event.preventDefault();
    }
    const response = await fetch(`http://localhost:8080/chats/${chatId}/messages`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'  
      },
      body: JSON.stringify({
        text: userMessage,
        task: {
          type: "goal_creation"
        }
      })
    });

    if (response.ok) {
      const responseData = await response.json();
      assistantMessage = responseData['text']
      userMessage = ""
    }
  }

  onMount(startNewChat)
</script>

<div class="container">
  <div class="spacer-1"></div>
  <div class="assistant-response-wrapper">
    <p id="assistant-response">{assistantMessage}</p>
  </div>
  <div class="spacer-2"></div>
  <form id="user-response">
    <input type="text" placeholder="Type your response here!" bind:value={userMessage}>
    <button on:click={fetchAssistantResponse}>-></button>
  </form>
  <div class="spacer-2"></div>
  <button on:click={startNewChat}>Reset</button>
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
