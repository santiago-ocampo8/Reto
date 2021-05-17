<template>
  <div>
      <h1>List of Buyers</h1>
      <buyers msg="Welcome to Your Vue.js App"/>
   </div>
</template>

<script>
// @ is an alias to /src
import buyers from '@/components/buyers.vue'
export default {
  name: 'Buyers',
  components: {
    buyers
  }
}
</script>
