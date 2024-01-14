class SearchApi {
  async fetchFastSearchList(value: string) {
    return await fetch(`/api/game/search/fast?searchKey=${value}`, {
      method: "GET",
    }).then((res) => res.text());
  }
}

export default new SearchApi();
