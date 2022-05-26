<script lang="ts">
    let email = "";
    let password = ""
</script>

<main>
    <form class="content">

        <p>Email:</p>
        <input type="email" required bind:value={email}>

        <p>Password:</p>
        <input type="password" required bind:value={password}>

        <br>
        <br>

        <button type="submit">
            Login
        </button>

        <button type="submit">
            Forgot password
        </button>
    </form>
</main>