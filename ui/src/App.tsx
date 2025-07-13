import {MagnifyingGlassIcon, PlusCircledIcon} from "@radix-ui/react-icons";
import {Box, Button, Flex, Text, TextField} from "@radix-ui/themes";
import {useState} from "react";
import {createParty} from "./api/partyClient";

function App() {
  const [partyName, setPartyName] = useState("");
  const [party, setParty] = useState<{ id: string; password: string } | null>(
    null,
  );

  const handleCreateParty = async () => {
    const result = await createParty(partyName);
    setParty(result);
  };

  return (
    <Flex gap={"4"} direction={"column"}>
      <Flex gap={"4"} direction={"row"} justify={"center"}>
        <Box width="200" height="16">
          <TextField.Root
            id="party-name"
            size={"2"}
            placeholder="Enter party name ..."
            value={partyName}
            onChange={(e) => setPartyName(e.target.value)}
            onKeyUp={(e) => {
              if (e.key === "Enter") {
                (async () => {
                  await handleCreateParty();
                })();
              }
            }}
          >
            <TextField.Slot>
              <MagnifyingGlassIcon height="16" width="16" />
            </TextField.Slot>
          </TextField.Root>
        </Box>
        <Box width={"64"} height={"16"}>
          <Button onClick={handleCreateParty}>
            <PlusCircledIcon />
            Create Party
          </Button>
        </Box>
      </Flex>
      {party && (
        <Flex gap={"4"} direction={"column"} align={"center"}>
          <Box>
            <Text as={"p"}>ID: {party.id}</Text>
          </Box>
          <Box>
            <Text as={"p"}>Password: {party.password}</Text>
          </Box>
        </Flex>
      )}
    </Flex>
  );
}

export default App;
