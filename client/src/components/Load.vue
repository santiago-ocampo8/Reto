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
        <v-btn
          align="center"
          color="red darken-4"
          dark
          @click="getDate"
          id="buttonLoad"
          :loading="loading"
        >
          Load
        </v-btn>
      </v-row>

      <v-row align="center" justify="space-around">
        <h3>{{ message }}</h3>
      </v-row>
    </v-container>
  </v-form>
</template>

<script>
export default {
  name: "Load",
  data: () => ({
    dateInput: "",
    message: "",
    loading: false,
  }),
  methods: {
    getDate() {
      var date = this.dateInput.toString();
      date = date.split("-");

      if (date[0] != "") {
    
        var dateFinal = (
          new Date(
            date[0],
            parseInt(date[1]) - 1,
            date[2]
          ).getTime() / 1000
        ).toFixed(0);
        this.loading = true;
        const url = "http://localhost:3000/" + "buyers/" + dateFinal;
        this.axios.post(url).then((response) => {
          console.log(response)
          this.message = "Date uploaded: "+date;
          this.loading = false;
        });
      } else {
        this.message = "Incorrect date";
      }
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
</style>
