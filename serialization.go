package neat

import (
    "encoding/json"
    "os"
    "io/ioutil"
)

func Save(neat *NEAT, path string) error {
    bytes, err := json.Marshal(neat.GetJsonNEAT())
    if err != nil {
        return err
    }

    err = ioutil.WriteFile(path, bytes, 0644)
    if err != nil {
        return err
    }

    return nil
}

func Load(path string) (*NEAT, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }

    decoder := json.NewDecoder(file)
    if err != nil {
        return nil, err
    }

    var jsonNeat JsonNEAT
    err = decoder.Decode(&jsonNeat)
    if err != nil {
        return nil, err
    }

    return MakeNEATFromJsonNEAT(jsonNeat), nil
}
