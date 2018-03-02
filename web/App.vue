<template>
  <div id="app" class="page">
    <!-- <router-view class="track-context" name="context"></router-view> -->
    <track-context class="track-context"></track-context>
    <router-view class="track-primary" role="main"></router-view>
    <!-- <router-view class="track-secondary" name="secondary"></router-view> -->
    <track-secondary class="track-secondary"></track-secondary>
    <footer class="track-footer"></footer>
  </div>
</template>

<script>
export default {
  name: 'app'
}
</script>
