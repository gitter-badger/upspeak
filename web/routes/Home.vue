<template>
  <div>
    <h1>Upspeak</h1>
    <p> A contextual discussion and knowledge-sharing platform</p>
    <button class="button" @click="doLogin">Login</button>
  </div>
</template>

<script>

export default {
  name: 'Home',
  methods: {
    doLogin: function () {
      console.log('Login')
    }
  }
}
</script>
