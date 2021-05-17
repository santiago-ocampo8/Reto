<template>
  <v-container>
    <v-row class="text-center">
      <v-col cols="12">
        <v-img
          :src="require('../assets/logo.png')"
          class="my-3"
          contain
          height="200"
        />
      </v-col>

      <v-col class="mb-4">
        <h1 class="display-2 font-weight-bold mb-3">
          Welcome To Report Shop
        </h1>

        <p class="subheading font-weight-regular">
          Santiago Ocampo Santana
          <br>Computer and Systems Engineering Student
          <br>Caldas University        
        </p>
      </v-col>

      
    </v-row>
  </v-container>
</template>

<script>
  export default {
    name: 'HelloWorld',
  }
</script>
