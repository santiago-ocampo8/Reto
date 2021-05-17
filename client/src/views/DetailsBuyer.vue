<template>
  <div>
      <h1>List of Buyers</h1>
      <details-Buyer />
   </div>
</template>

<script>
// @ is an alias to /src
import detailsBuyer from '@/components/details.vue'
export default {
  name: 'DetailsBuyer',
  components: {
      detailsBuyer
  }
}
</script>
