<template>
  <v-form>
    <v-container>
      <v-row align="center" justify="space-around">
        <v-col cols="12" sm="6" md="3">
          <v-text-field
            type="date"
            label="Date"
            v-model="dateInput"
          ></v-text-field>
        </v-col>
      </v-row>
      <v-row align="center" justify="space-around">
        <v-btn depressed color="red darken-4" dark @click="getDate">
          Load
        </v-btn>
      
      </v-row>
      <h3>{{message}}</h3>
    </v-container>
  </v-form>
</template>

<script>
export default {
  name: "Load",
  data: () => ({
    dateInput: "",
    message:""

  }),
  methods: {
    getDate() {
      var temporal = this.dateInput.toString();
      temporal = temporal.split("-");
      if (temporal[0]!=""){
        document.getElementById("buttonBUenas").style.display = "none";
        console.log(typeof this.dateInput);
        console.log(temporal);
      var hola = (
        new Date(
          temporal[0],
          parseInt(temporal[1]) - 1,
          temporal[2]
        ).getTime() / 1000
      ).toFixed(0);
      const url = "http://localhost:3000/" + "buyers/" + hola;
      this.axios.post(url).then((response) => {
        console.log(response);
        document.getElementById("buttonBUenas").style.display = "block";
        this.message="Goog"
      });
      }else{
          this.message="Bad"
      }
    }
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
</style>
